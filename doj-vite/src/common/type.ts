export type TagType={id:number,name:string,color:string}

export type ProblemListType={
    problem_id:string,
    title:string,
    difficulty:number,
    myStatus:number,
    tags:TagType[],
    total:number,
    ac:number,
    wa:number,
    tle:number,
    mle:number,
    re:number,
    pe:number,
    ce:number,
    se:number
}
