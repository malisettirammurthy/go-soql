/*
 * Copyright (c) 2018, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package soql_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/forcedotcom/go-soql"
)

var _ = Describe("Marshaller", func() {
	Describe("MarshalWhereClause", func() {
		var (
			clause         string
			expectedClause string
			err            error
		)
		Context("when non pointer value is passed as argument", func() {
			var (
				critetria TestQueryCriteria
			)

			JustBeforeEach(func() {
				clause, err = MarshalWhereClause(critetria)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when there are no fields populated", func() {
				It("returns empty where clause", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when only like clause pattern is populated", func() {
				Context("when there is only one item in the like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							IncludeNamePattern: []string{"-db"},
						}
						expectedClause = "Host_Name__c LIKE '%-db%'"
					})

					It("returns where clause with only one condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							IncludeNamePattern: []string{"-db", "-dbmgmt", "-dgdb"},
						}
						expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%' OR Host_Name__c LIKE '%-dgdb%')"
					})

					It("returns where clause with OR condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only not like clause is populated", func() {
				Context("when there is only one item in the not like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							ExcludeNamePattern: []string{"-db"},
						}
						expectedClause = "(NOT Host_Name__c LIKE '%-db%')"
					})

					It("returns where clause with only one condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the not like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							ExcludeNamePattern: []string{"-db", "-dbmgmt", "-dgdb"},
						}
						expectedClause = "((NOT Host_Name__c LIKE '%-db%') AND (NOT Host_Name__c LIKE '%-dbmgmt%') AND (NOT Host_Name__c LIKE '%-dgdb%'))"
					})

					It("returns where clause with OR condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only equalsClause is populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						AssetType: "SERVER",
					}
					expectedClause = "Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns appropriate where clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when only notEqualsClause is populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						Status: "InActive",
					}
					expectedClause = "Status__c != 'InActive'"
				})

				It("returns appropriate where clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when only inClause is populated", func() {
				Context("when there is only one item in the inClause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							Roles: []string{"db"},
						}
						expectedClause = "Role__r.Name IN ('db')"
					})
					It("returns where clause with only one item in IN clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the inClause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							Roles: []string{"db", "dbmgmt"},
						}
						expectedClause = "Role__r.Name IN ('db','dbmgmt')"
					})
					It("returns where clause with all the items in IN clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only null clause is populated", func() {
				Context("when null is allowed", func() {
					BeforeEach(func() {
						allowNull := true
						critetria = TestQueryCriteria{
							AllowNullLastDiscoveredDate: &allowNull,
						}

						expectedClause = "Last_Discovered_Date__c = null"
					})

					It("returns appropriate where clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when null is not allowed", func() {
					BeforeEach(func() {
						allowNull := false
						critetria = TestQueryCriteria{
							AllowNullLastDiscoveredDate: &allowNull,
						}

						expectedClause = "Last_Discovered_Date__c != null"
					})

					It("returns appropriate where clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when likeClause and inClause are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db')"
				})

				It("returns properly formed clause for name and role joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when likeClause and equalsClause are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						AssetType:          "SERVER",
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns properly formed clause for likeClause and inClause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when both likeClause and notLikeClause are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						ExcludeNamePattern: []string{"-core", "-drp"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%'))"
				})

				It("returns properly formed clause for likeClause and notLikeClause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when all clauses are populated", func() {
				BeforeEach(func() {
					allowNull := false
					critetria = TestQueryCriteria{
						AssetType:                   "SERVER",
						IncludeNamePattern:          []string{"-db", "-dbmgmt"},
						Roles:                       []string{"db", "dbmgmt"},
						ExcludeNamePattern:          []string{"-core", "-drp"},
						AllowNullLastDiscoveredDate: &allowNull,
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%')) AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER' AND Last_Discovered_Date__c != null"
				})

				It("returns properly formed clause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when struct has invalid tag key", func() {
				type InvalidCriteriaStruct struct {
					SomePattern      []string `soql:"likeClause,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"invalidClause,fieldName=Some_Other_Field"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalWhereClause(InvalidCriteriaStruct{})
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has missing fieldName", func() {
				type MissingFieldName struct {
					SomePattern      []string `soql:"likeClause,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"equalsClause"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalWhereClause(MissingFieldName{})
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has invalid types", func() {
				It("returns empty string", func() {
					str, err := MarshalWhereClause(QueryCriteriaWithInvalidTypes{})
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeEmpty())
				})
			})
		})

		Context("when pointer is passed as argument", func() {
			var (
				critetria *TestQueryCriteria
			)

			JustBeforeEach(func() {
				clause, err = MarshalWhereClause(critetria)
			})

			Context("when nil is passed as argument", func() {
				It("returns empty where clause", func() {
					Expect(err).To(Equal(ErrNilValue))
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when empty value is passed as argument", func() {
				BeforeEach(func() {
					critetria = &TestQueryCriteria{}
				})

				It("returns empty where clause", func() {
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when all values are populated", func() {
				BeforeEach(func() {
					critetria = &TestQueryCriteria{
						AssetType:          "SERVER",
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
						ExcludeNamePattern: []string{"-core", "-drp"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%')) AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns properly formed clause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})
		})
	})

	Describe("MarshalSelectClause", func() {
		Context("when non pointer value is passed as argument", func() {
			Context("when no relationship name is passed", func() {
				Context("when no nested struct is passed", func() {
					It("returns just the json tag names of fields concatenanted by comma", func() {
						str, err := MarshalSelectClause(NonNestedStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Name,SomeValue__c"))
					})
				})

				Context("when nested struct is passed", func() {
					It("returns properly resolved list of field names", func() {
						str, err := MarshalSelectClause(NestedStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c"))
					})
				})
			})

			Context("when relationship name is passed", func() {
				Context("when no nested struct is passed", func() {
					It("returns just the json tag names of fields concatenanted by comma and prefixed by relationship name", func() {
						str, err := MarshalSelectClause(NonNestedStruct{}, "Role__r")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Role__r.Name,Role__r.SomeValue__c"))
					})
				})
			})

			Context("when struct has invalid tag key", func() {
				type InvalidStruct struct {
					Id  string `soql:"selectColumn,fieldName=Id"`
					Foo string `soql:"invalidClause,fieldName=Foo"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalSelectClause(InvalidStruct{}, "")
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has missing fieldName", func() {
				type MissingFieldName struct {
					SomePattern      []string `soql:"selectColumn,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"selectColumn"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalSelectClause(MissingFieldName{}, "")
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has child relationship", func() {
				Context("when child struct has select clause only", func() {
					It("returns properly constructed select clause", func() {
						str, err := MarshalSelectClause(ParentStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Application_Versions__c.Version__c FROM Application_Versions__r)"))
					})
				})

				Context("when child struct has select clause and where clause", func() {
					It("returns properly constructed select clause", func() {
						str, err := MarshalSelectClause(ParentStruct{
							ChildStruct: TestChildStruct{
								WhereClause: ChildQueryCriteria{
									Name: "sfdc-release",
								},
							},
						}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Application_Versions__c.Version__c FROM Application_Versions__r WHERE SM_Application_Versions__c.Name__c = 'sfdc-release')"))
					})
				})

				Context("when child struct does not have select clause", func() {
					It("returns error", func() {
						_, err := MarshalSelectClause(InvalidParentStruct{}, "")
						Expect(err).To(Equal(ErrNoSelectClause))
					})
				})

				Context("when selectChild tag is applied to non struct member", func() {
					It("returns error", func() {
						_, err := MarshalSelectClause(ChildTagToNonStruct{}, "")
						Expect(err).To(Equal(ErrInvalidTag))
					})
				})
			})
		})

		Context("when pointer value is passed as argument", func() {
			Context("when nil is passed", func() {
				It("returns ErrNilValue error", func() {
					var r *NestedStruct
					str, err := MarshalSelectClause(r, "")
					Expect(err).To(Equal(ErrNilValue))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when nested struct is passed", func() {
				It("returns properly resolved list of field names", func() {
					str, err := MarshalSelectClause(&NestedStruct{}, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c"))
				})
			})
		})
	})

	Describe("Marshal", func() {
		var (
			soqlStruct    interface{}
			expectedQuery string
			actualQuery   string
			err           error
		)

		JustBeforeEach(func() {
			actualQuery, err = Marshal(soqlStruct)
		})

		Context("when empty struct is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = EmptyStruct{}
			})

			It("returns empty string", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(BeEmpty())
			})
		})

		Context("when valid value is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = TestSoqlStruct{
					SelectClause: NestedStruct{},
					WhereClause: TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')"
			})

			It("returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})

		Context("when valid pointer is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = &TestSoqlStruct{
					SelectClause: NestedStruct{},
					WhereClause: TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')"
			})

			It("returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})

		Context("when struct with no soql tags is passed", func() {
			BeforeEach(func() {
				soqlStruct = NonSoqlStruct{}
			})

			It("returns emptyString", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(BeEmpty())
			})
		})

		Context("when struct with multiple selectClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleSelectClause{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrMultipleSelectClause))
			})
		})

		Context("when struct with multiple whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleWhereClause{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrMultipleWhereClause))
			})
		})

		Context("when struct with only whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = OnlyWhereClause{
					WhereClause: TestQueryCriteria{
						AssetType:          "SERVER",
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrNoSelectClause))
			})
		})

		Context("when struct with multiple whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleWhereClause{
					WhereClause1: ChildQueryCriteria{
						Name: "foo",
					},
					WhereClause2: ChildQueryCriteria{
						Name: "bar",
					},
				}
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when nil pointer is passed", func() {
			BeforeEach(func() {
				var ptr *TestSoqlStruct
				soqlStruct = ptr
			})

			It("returns ErrNilValue error", func() {
				Expect(err).To(Equal(ErrNilValue))
			})
		})

		Context("when struct with invalid tag is passed", func() {
			BeforeEach(func() {
				soqlStruct = InvalidTagInStruct{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrInvalidTag))
			})
		})
	})
})